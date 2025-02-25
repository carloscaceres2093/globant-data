import os
import math
import boto3
import psycopg2
import gzip
import io
import logging
import logging.config
import yaml
from concurrent.futures import ThreadPoolExecutor, as_completed
from threading import Lock
import time
import fastavro
import csv
from dotenv import load_dotenv
load_dotenv()

# AWS_ACCESS_KEY = os.getenv("AWS_ACCESS_KEY_ID")
# AWS_SECRET_KEY = os.getenv("AWS_SECRET_ACCESS_KEY")
# REGION_NAME = os.getenv("REGION_NAME")
TABLE_NAME = os.getenv("TABLE_NAME")

# S3_BUCKET = os.getenv("LANDING_BUCKET")
DB_NAME = os.getenv("DB_NAME")
NAME_PART = os.getenv("NAME_PART")
name_part = NAME_PART if NAME_PART else ""
S3_CHECKPOINT_PATH = f"data/backup/checkpoints/"
checkpoint_file_name = "checkpoint_{DB_NAME}_{TABLE_NAME}.txt"
S3_DUMP_PATH_TEMPLATE = "data/backup/{db_name}{table}{name_add}/"


DB_CONFIG = {
    "host": os.getenv("DB_HOST"),
    "port": 5432,
    "dbname": DB_NAME,
    "user": os.getenv("DB_USER"),
    "password": os.getenv("DB_PASSWORD"),
}
print(DB_CONFIG)
ROWS_PER_DUMP = int(os.getenv("ROWS_PER_DUMP"))
# s3_client = boto3.client(
#     "s3",
#     aws_access_key_id=AWS_ACCESS_KEY,
#     aws_secret_access_key=AWS_SECRET_KEY,
#     region_name=REGION_NAME,
# )
lock = Lock()
processed_parts = set()
failed_parts = set()


def get_db_connection():
    return psycopg2.connect(**DB_CONFIG)


def get_row_count():
    with get_db_connection() as conn:
        with conn.cursor() as cur:
            cur.execute(
                f"SELECT reltuples::BIGINT AS estimated_count FROM pg_class WHERE relname = '{TABLE_NAME}';"
            )
            return cur.fetchone()[0]


def get_min_max_ids(codition=""):
    with get_db_connection() as conn:
        with conn.cursor() as cur:
            cur.execute(f"SELECT MIN(id), MAX(id) FROM {TABLE_NAME} {codition}")
            return cur.fetchone()


# def get_checkpoint():
#     try:
#         response = s3_client.get_object(Bucket=S3_BUCKET, Key=S3_CHECKPOINT_PATH)
#         checkpoint = response["Body"].read().decode("utf-8")
#         if checkpoint == "finished":
#             return 0
#         return int(checkpoint)
#     except s3_client.exceptions.NoSuchKey:
#         return 0

def get_checkpoint():

    # Define the path

    # Check if the path exists
    if not os.path.exists(S3_CHECKPOINT_PATH):
        # Create the directory (including parent directories if needed)
        os.makedirs(S3_CHECKPOINT_PATH)
        print(f"Directory created: {S3_CHECKPOINT_PATH}")
        return 0
    else:
        with open(S3_CHECKPOINT_PATH+checkpoint_file_name, mode="r", encoding="utf-8") as file:
            checkpoint = file.read()

        # Now `checkpoint` contains the content of the file
        print(checkpoint)
        if checkpoint == "finished":
            return 0
        return int(checkpoint)

# def save_checkpoint(part):
#     s3_client.put_object(Bucket=S3_BUCKET, Key=S3_CHECKPOINT_PATH, Body=str(part))

def infer_avro_type(value):
    """
    Infers the Avro type for a given value.
    """
    try:
        int(value)
        return "int"
    except ValueError:
        try:
            float(value)
            return "double"
        except ValueError:
            if value.lower() in ("true", "false"):
                return "boolean"
            else:
                return "string"

def generate_avro_schema(fields, sample_records, table_name):
    """
    Generates an Avro schema dynamically based on the CSV header and sample records.
    """
    schema = {
        "type": "record",
        "name": table_name,
        "fields": [],
    }

    for field in fields:
        # Infer the type based on the first non-null value in the sample records
        field_type = None
        for record in sample_records:
            if record[field]:  # Skip null/empty values
                field_type = "string"#infer_avro_type(record[field])
                break
        if not field_type:
            field_type = "string"  # Default to string if no data is available

        schema["fields"].append({"name": field, "type": field_type})
    return schema


def dump_batch(part, start_id, end_id):
    dump_path = S3_DUMP_PATH_TEMPLATE.format(
        table=TABLE_NAME, db_name=DB_NAME, name_add=name_part
    )
    dump_file_name = f"dump_part_{part}.avro"
    query = f"COPY (SELECT * FROM {TABLE_NAME} WHERE id >= {start_id} AND id < {end_id}) TO STDOUT WITH CSV HEADER"
    logger.info(f"Building query with start: {start_id} and limit id: {end_id}")
    with get_db_connection() as connection:
        with connection.cursor() as cur:
            csv_buffer = io.StringIO()
            cur.copy_expert(query, csv_buffer)
            csv_buffer.seek(0)

            # Read the CSV data into a list of dictionaries
            csv_reader = csv.DictReader(csv_buffer)
            records = [row for row in csv_reader]
            if not records:
                logger.warning("No records found for the given query.")
                return
            # Generate the Avro schema dynamically
            schema = generate_avro_schema(csv_reader.fieldnames, records, TABLE_NAME)
            # Write the records to an Avro file
            if not os.path.exists(dump_path):
                # Create the directory (including parent directories if needed)
                os.makedirs(dump_path)
                print(f"Directory created: {dump_path}")
            with open(dump_path+dump_file_name, "wb") as avro_file:
                fastavro.writer(avro_file, schema, records)



def parallel_dump(start_part, min_id, max_id, rows_per_dump, max_workers):
    def task(part, start_id, end_id):
        retries = 5
        logger.info(f"Starting part {part}: IDs {start_id} to {end_id}")
        with lock:
            if part in processed_parts:
                logger.info(f"Skipping part {part}, already in progress")
                return
            processed_parts.add(part)

        try:
            for attempt in range(1, retries + 1):
                try:
                    dump_batch(part, start_id, end_id)
                    logger.info(f"Completed dump part {part}")
                    #save_checkpoint(part)
                    break
                except Exception as e:
                    logger.error(
                        f"Error in dump part {part}, attempt {attempt}/{retries}: {e}"
                    )
                    if attempt == retries:
                        logger.critical(
                            f"Max retries reached for part {part}. Adding to failed_parts."
                        )
                        with lock:
                            failed_parts.add(part)
                    time.sleep(5 * attempt)  # Exponential backoff
        finally:
            with lock:
                processed_parts.discard(part)

    parts = [
        (
            part,
            min_id + (part - 1) * rows_per_dump,
            min(min_id + part * rows_per_dump, max_id + 1),
        )
        for part in range(
            start_part, math.ceil((max_id - min_id + 1) / rows_per_dump) + 1
        )
    ]
    for part, start_id, end_id in parts:
        task(part, start_id, end_id)
    # with ThreadPoolExecutor(max_workers=max_workers) as executor:
    #     futures = {
    #         executor.submit(task, part, start_id, end_id): part
    #         for part, start_id, end_id in parts
    #     }
    #     for future in as_completed(futures):
    #         part = futures[future]
    #         try:
    #             future.result()
    #         except Exception as e:
    #             logger.error(f"Task failed with exception in part {part}: {e}")

    # Retry failed parts
    if failed_parts:
        logger.warning(f"Retrying failed parts: {failed_parts}")
        for part in list(failed_parts):
            start_id = min_id + (part - 1) * rows_per_dump
            end_id = min(start_id + rows_per_dump, max_id + 1)
            try:
                task(part, start_id, end_id)
                with lock:
                    failed_parts.discard(part)
            except Exception as e:
                logger.error(f"Part {part} failed again: {e}")


if __name__ == "__main__":
    with open("/app/tasks/config/logging.yaml", "r") as f:
        config = yaml.safe_load(f.read())
        logging.config.dictConfig(config)
    logger = logging.getLogger(__name__)

    total_rows = get_row_count()
    min_id, max_id = get_min_max_ids()
    print(min_id, max_id)
    
    checkpoint = get_checkpoint()
    
    logger.info(f"Total rows (estimated): {total_rows}")
    logger.info(f"Min ID: {min_id}, Max ID: {max_id}")
    logger.info(f"Rows per dump: {ROWS_PER_DUMP}")
    logger.info(f"Resuming from checkpoint: {checkpoint}")

    parallel_dump(checkpoint + 1, min_id, max_id, ROWS_PER_DUMP, max_workers=4)
    quit()
    save_checkpoint("finished")