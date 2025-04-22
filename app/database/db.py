import psycopg2
from psycopg2 import sql

# Database connection parameters
db_config = {
    "dbname": "lvg",
    "user": "dbuser",
    "password": "dbpass",
    "host": "localhost",
    "port": "5432"
}

# Connect to the database
try:
    # Establish a connection to PostgreSQL
    conn = psycopg2.connect(**db_config)
    conn.autocommit = True  # Automatically commit transactions (for simplicity)
    cursor = conn.cursor()

    # Define SQL statement for creating a new table
    create_table_query = sql.SQL("""
        CREATE TABLE IF NOT EXISTS lvg_requests (
            job_id UUID PRIMARY KEY,
            timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            values JSONB
        )
    """)

    # Execute the SQL statement
    cursor.execute(create_table_query)
    print("Table created successfully!")

except Exception as e:
    print(f"An error occurred: {e}")
finally:
    # Clean up and close the connection
    if cursor:
        cursor.close()
    if conn:
        conn.close()
