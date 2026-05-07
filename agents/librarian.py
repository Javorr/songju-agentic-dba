import psycopg2
from typing import Dict, List, Optional

from .config import Config
from .sql_utils import extract_tables


class LibrarianAgent:
    """
    Agent A: Identifies relevant table schemas from query metadata.
    MVP: Simple lookup from PostgreSQL metadata store.
    """

    def __init__(self, config: Config):
        self.config = config
        self.conn = None

    def connect(self):
        """Connect to PostgreSQL metadata store."""
        self.conn = psycopg2.connect(self.config.database_url)
        print("Librarian connected to PostgreSQL")

    def get_table_schema(self, table_name: str) -> Optional[Dict]:
        """
        Fetch table schema from metadata store.
        Returns: {'table': 'users', 'columns': [{'name': 'id', 'type': 'int'}, ...]}
        """
        if not self.conn:
            self.connect()

        with self.conn.cursor() as cur:
            cur.execute("""
                SELECT column_name, data_type
                FROM information_schema.columns
                WHERE table_name = %s
                ORDER BY ordinal_position
            """, (table_name,))

            columns = [{'name': row[0], 'type': row[1]} for row in cur.fetchall()]

            if columns:
                return {'table': table_name, 'columns': columns}
            return None

    def analyze_query(self, sql: str, service: str) -> Dict:
        """
        MVP: Extract tables from SQL and fetch their schemas.
        """
        tables = extract_tables(sql)
        schemas = []

        for table in tables:
            schema = self.get_table_schema(table)
            if schema:
                schemas.append(schema)

        return {
            'service': service,
            'tables_found': tables,
            'schemas': schemas
        }
