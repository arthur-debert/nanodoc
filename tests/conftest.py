import logging
import os
import sys

import pytest

# Add the src directory to the Python path
sys.path.insert(
    0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "src"))
)

# Set verbose logging for all tests
os.environ["NANODOC_VERBOSE"] = "1"

# Configure root logger for tests
logging.basicConfig(
    level=logging.DEBUG,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)


@pytest.fixture(autouse=True)
def setup_logging():
    """Configure logging for all tests automatically."""
    # Configure all module loggers
    module_names = [
        "nanodoc",  # Main logger
        "cli",  # V2 CLI
        "document",  # V2 Document
        "formatter",  # V2 Formatter
        "renderer",  # V2 Renderer
        "resolver",  # V2 Resolver
        "extractor",  # V2 Extractor
    ]

    for name in module_names:
        module_logger = logging.getLogger(name)
        module_logger.setLevel(logging.DEBUG)
        module_logger.propagate = False

        # Remove any existing handlers to prevent duplicate logging
        module_logger.handlers.clear()

        # Add handler if the logger doesn't have one
        if not module_logger.handlers:
            handler = logging.StreamHandler(sys.stdout)
            fmt = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
            formatter = logging.Formatter(fmt)
            handler.setFormatter(formatter)
            module_logger.addHandler(handler)

    yield

    # Clean up after test
    for name in module_names:
        logging.getLogger(name).handlers.clear()


@pytest.fixture
def caplog(caplog):
    """Fixture to capture log messages."""
    caplog.set_level(logging.DEBUG)
    return caplog
