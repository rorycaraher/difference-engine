import uuid
import psycopg2
from datetime import datetime

# Connect to your PostgreSQL database
conn = psycopg2.connect(
    dbname="de",
    user="dbuser",
    password="dbpass",
    host="localhost",
    port="5432"
)
cursor = conn.cursor()

# Generate a UUID
job_id = str(uuid.uuid4())
timestamp = datetime.now()

# Insert into table
cursor.execute("""
    INSERT INTO de_requests (job_id, timestamp, values) VALUES (%s, %s, %s)
""", (job_id, timestamp , '{"key": "value"}'))

conn.commit()
cursor.close()
conn.close()
