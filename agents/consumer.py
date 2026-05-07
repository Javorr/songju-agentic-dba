import asyncio
import signal
from typing import Optional

import nats
from nats.js.api import ConsumerConfig, AckPolicy

from .config import Config
from .librarian import LibrarianAgent


class SlowLogsConsumer:
    """Consumes slow query logs from NATS JetStream."""

    def __init__(self, config: Config):
        self.config = config
        self.nc = None
        self.js = None
        self.subscription = None
        self.librarian = LibrarianAgent(config)
        self.running = False

    async def connect(self):
        """Connect to NATS and setup JetStream."""
        self.nc = await nats.connect(self.config.NATS_URL)
        self.js = self.nc.jetstream()
        self.librarian.connect()
        print("Connected to NATS and PostgreSQL")

    async def setup_consumer(self, durable_name: str = "librarian-processor"):
        """Setup pull consumer for db.logs.slow."""
        config = ConsumerConfig(
            durable_name=durable_name,
            filter_subject="db.logs.slow",
            ack_policy=AckPolicy.EXPLICIT,
            max_deliver=3,
            ack_wait=30.0,
        )

        self.subscription = await self.js.pull_subscribe(
            subject="db.logs.slow",
            durable=durable_name,
            stream="DB_LOGS",
            config=config
        )
        print(f"Consumer '{durable_name}' ready")

    async def handle_message(self, msg):
        """Process individual message."""
        try:
            import json
            payload = json.loads(msg.data.decode())
            sql = payload.get('sql', '')
            service = payload.get('service', 'unknown')

            print(f"Processing: {sql[:80]}...")

            # Librarian agent analyzes the query
            result = self.librarian.analyze_query(sql, service)
            print(f"Librarian result: {result}")

            await msg.ack()
            print("Message acknowledged")

        except json.JSONDecodeError:
            print("Invalid JSON, terminating message")
            await msg.term()
        except Exception as e:
            print(f"Error: {e}")
            await msg.nak()

    async def process_messages(self, batch_size: int = 10):
        """Main processing loop."""
        self.running = True

        while self.running:
            try:
                msgs = await self.subscription.fetch(batch_size, timeout=5.0)
                for msg in msgs:
                    await self.handle_message(msg)
            except nats.errors.TimeoutError:
                await asyncio.sleep(1)
            except Exception as e:
                print(f"Fetch error: {e}")
                await asyncio.sleep(2)

    async def close(self):
        """Clean shutdown."""
        self.running = False
        if self.librarian.conn:
            self.librarian.conn.close()
        if self.nc:
            await self.nc.close()
