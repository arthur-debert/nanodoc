"""Test configuration module."""

import logging
import os
import sys

import pytest

from nanodoc.v2.boot import MODULE_LOGGERS

# Add the src directory to the Python path
sys.path.insert(
    0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "src"))
)

# Set verbose logging for all tests
os.environ["NANODOC_VERBOSE"] = "1"


def configure_test_logging(loggers):
    """Configure logging for tests.

    Args:
        loggers: List of logger names to configure.
    """
    for name in loggers:
        logger = logging.getLogger(name)
        logger.setLevel(logging.DEBUG)
        logger.handlers.clear()
        handler = logging.StreamHandler(sys.stdout)
        handler.setFormatter(logging.Formatter("%(message)s"))
        logger.addHandler(handler)


def cleanup_test_logging(loggers):
    """Clean up logging configuration after tests.

    Args:
        loggers: List of logger names to clean up.
    """
    for name in loggers:
        logger = logging.getLogger(name)
        logger.handlers.clear()


@pytest.fixture(autouse=True)
def setup_logging():
    """Configure logging for all tests automatically."""
    configure_test_logging(MODULE_LOGGERS)
    yield
    cleanup_test_logging(MODULE_LOGGERS)


@pytest.fixture
def caplog(caplog):
    """Fixture to capture log messages."""
    caplog.set_level(logging.DEBUG)
    return caplog
