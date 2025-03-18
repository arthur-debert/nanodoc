"""Test configuration module."""

import logging
import os
import sys

import pytest

from tests.boot import (
    MODULE_LOGGERS,
    cleanup_test_logging,
    configure_test_logging,
)

# Add the src directory to the Python path
sys.path.insert(
    0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "src"))
)

# Set verbose logging for all tests
os.environ["NANODOC_VERBOSE"] = "1"


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
