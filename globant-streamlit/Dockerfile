# Use Python as base image
FROM python:3.10

# Set working directory
WORKDIR /app

# Copy application files
COPY streamlit_app.py /app

# Install dependencies
RUN pip install streamlit requests pandas matplotlib seaborn

# Expose Streamlit's default port
EXPOSE 8501

# Run the Streamlit app
CMD ["streamlit", "run", "streamlit_app.py", "--server.port=8501", "--server.address=0.0.0.0"]