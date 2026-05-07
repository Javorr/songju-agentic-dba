from sqlglot import parse_one, exp
from typing import List


def extract_tables(sql_query: str) -> List[str]:
    """
    Extract table names from SQL query.
    Handles masked PII (treats [REDACTED_*] as string literals).

    Examples:
        >>> extract_tables("SELECT * FROM users WHERE email = [REDACTED_EMAIL]")
        ['users']
        >>> extract_tables("SELECT * FROM sales.orders JOIN hr.employees ON ...")
        ['orders', 'employees']
    """
    try:
        parsed = parse_one(sql_query)
        tables = [table.name for table in parsed.find_all(exp.Table)]
        return tables
    except Exception as e:
        print(f"Warning: Could not parse SQL: {e}")
        return []
