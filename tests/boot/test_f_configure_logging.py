"""Tests for boot.configure_logging function."""

import logging
from unittest.mock import patch

from nanodoc.boot import configure_logging


def test_configure_logging_default():
    """Test configure_logging with default settings."""
    # Reset logging before test
    root = logging.getLogger()
    for handler in root.handlers[:]:
        root.removeHandler(handler)

    # Test with default settings
    with patch.object(logging, "basicConfig") as mock_basic_config:
        configure_logging()
        mock_basic_config.assert_called_once()
        args, kwargs = mock_basic_config.call_args
        assert kwargs["level"] == logging.WARNING


def test_configure_logging_verbose():
    """Test configure_logging with verbose flag."""
    # Reset logging before test
    root = logging.getLogger()
    for handler in root.handlers[:]:
        root.removeHandler(handler)

    # Test with verbose=True
    with patch.object(logging, "basicConfig") as mock_basic_config:
        configure_logging(verbose=True)
        mock_basic_config.assert_called_once()
        args, kwargs = mock_basic_config.call_args
        assert kwargs["level"] == logging.DEBUG


def test_configure_logging_env_var():
    """Test configure_logging with environment variable."""
    # Reset logging before test
    root = logging.getLogger()
    for handler in root.handlers[:]:
        root.removeHandler(handler)

    # Test with NANODOC_VERBOSE environment variable
    with (
        patch.object(logging, "basicConfig") as mock_basic_config,
        patch.dict("os.environ", {"NANODOC_VERBOSE": "1"}),
    ):
        configure_logging()
        mock_basic_config.assert_called_once()
        args, kwargs = mock_basic_config.call_args
        assert kwargs["level"] == logging.DEBUG
