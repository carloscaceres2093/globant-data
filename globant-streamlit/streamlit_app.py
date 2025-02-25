import streamlit as st
import requests
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import os

st.set_page_config(layout="wide")
st.title("📊 Company Metrics Dashboard")

col1, col2, col3 = st.columns(3)

with col1:
    selected_year = st.text_input("Enter Year", "2021")
with col2:
    department_name = st.text_input("Enter Department (Leave empty for all departments)", "")
with col3:
    job_name = st.text_input("Enter Job Name (Leave empty for all jobs)", "")

# Compose URL
base_path = os.getenv("BASE_PATH")
base_url_quarterly = f"{base_path}/quarter_metrics"
base_url_hired = f"{base_path}/hired_metrics"
query_params = []
if not selected_year:
    selected_year ="2021"
query_params.append(f"year={selected_year}")
if job_name.strip(): 
    query_params.append(f"job_name={job_name.replace(' ', '%20')}")
if department_name.strip(): 
    query_params.append(f"department_name={department_name.replace(' ', '%20')}")

api_url_quarterly = f"{base_url_quarterly}?{'&'.join(query_params)}" if query_params else base_url_quarterly
api_url_hired =  f"{base_url_hired}?{'&'.join(query_params)}" if query_params else base_url_hired

# Fetch Data
@st.cache_data
def get_data(api_url):
    response = requests.get(api_url)
    if response.status_code == 200:
        return response.json()
    else:
        st.error(f"Failed to fetch data from {api_url}")
        return []

data_quarterly = get_data(api_url_quarterly)
data_hired = get_data(api_url_hired)


# CHARTS
quarter_filter = st.selectbox("Select a Quarter:", ["All", "q_1", "q_2", "q_3", "q_4"])

# QUARTER CHARTS
if data_quarterly:
    df_quarterly = pd.DataFrame(data_quarterly)

    df_melted = df_quarterly.melt(id_vars=["department"], value_vars=["q_1", "q_2", "q_3", "q_4"], var_name="Quarter", value_name="Count")

    department_options = ["All Departments"] + list(df_quarterly["department"].unique())
    selected_department = st.selectbox("Select a department:", department_options, key="quarterly")

    fig, ax = plt.subplots(figsize=(16, 6))

    if selected_department == "All Departments":
        for dept in df_quarterly["department"].unique():
            dept_data = df_melted[df_melted["department"] == dept]
            ax.plot(dept_data["Quarter"], dept_data["Count"], marker="o", linestyle="-", label=dept)
    else:
        filtered_df = df_melted[df_melted["department"] == selected_department]
        ax.bar(filtered_df["Quarter"], filtered_df["Count"], color='blue')

    ax.set_title(f"Quarterly Metrics - {selected_department if selected_department != 'All Departments' else 'All'}", fontsize=16)
    ax.set_xlabel("Quarter", fontsize=14)
    ax.set_ylabel("Count", fontsize=14)
    ax.legend(loc="upper right", fontsize=12)
    ax.grid(True, linestyle="--", alpha=0.5)

    st.pyplot(fig, use_container_width=True)
else:
    st.warning("⚠️ No quarterly metrics data available.")

if data_quarterly:
    fig, ax = plt.subplots(figsize=(16, 6))
    df_melted = df_quarterly.melt(id_vars=["department", "job"], value_vars=["q_1", "q_2", "q_3", "q_4"], 
                                   var_name="Quarter", value_name="Count")
    heatmap_data = df_melted.pivot_table(index="department", columns="job", values="Count", aggfunc="sum", fill_value=0)
    sns.heatmap(heatmap_data, annot=False, fmt="d", cmap="Blues", linewidths=0.5, ax=ax)

    ax.set_title("🔥 Job-Department Relationship (Heatmap)", fontsize=16)
    ax.set_xlabel("Job", fontsize=14)
    ax.set_ylabel("Department", fontsize=14)

    st.pyplot(fig, use_container_width=True)

# HIRED  CHART 
if data_hired:
    df_hired = pd.DataFrame(data_hired)

    df_hired = df_hired.sort_values(by="hired", ascending=False)
    fig, ax = plt.subplots(figsize=(16, 6))
    ax.barh(df_hired["department"], df_hired["hired"], color="green")

    ax.set_title(f"💼 Hired Employees by Department ({selected_year})", fontsize=16)
    ax.set_xlabel("Number of Hires", fontsize=14)
    ax.set_ylabel("Department", fontsize=14)
    ax.grid(axis="x", linestyle="--", alpha=0.5)

    st.pyplot(fig, use_container_width=True)
else:
    st.warning("⚠️ No hired metrics data available.")