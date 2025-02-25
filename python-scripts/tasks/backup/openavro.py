import fastavro

# Path to the Avro file
avro_file_path = "/home/carlos/Documents/go/python-scripts/data/backup/globant-msemployees/dump_part_1.avro"

# Open the Avro file and read its content
with open(avro_file_path, "rb") as avro_file:
    # Use fastavro.reader to read the Avro file
    reader = fastavro.reader(avro_file)
    
    # Print the schema of the Avro file
    print("Schema:", reader.writer_schema)
    
    # Iterate over each record in the Avro file
    for record in reader:
        print(record)