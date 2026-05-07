import asyncio
import signal
from .consumer import SlowLogsConsumer
from .config import Config


async def main():
    config = Config()
    consumer = SlowLogsConsumer(config)

    # Setup signal handlers
    loop = asyncio.get_event_loop()
    for sig in (signal.SIGTERM, signal.SIGINT):
        loop.add_signal_handler(sig, lambda: asyncio.create_task(consumer.close()))

    try:
        await consumer.connect()
        await consumer.setup_consumer()
        await consumer.process_messages()
    except KeyboardInterrupt:
        print("Interrupted by user")
    finally:
        await consumer.close()


if __name__ == '__main__':
    asyncio.run(main())
