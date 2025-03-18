"""Test configuration and setup module.

This module handles test environment configuration, including logging setup.
"""

import logging
import sys


def configure_test_logging(module_names: list[str]) -> None:
    """Configure logging for test environment.

    Args:
        module_names: List of module logger names to configure.
    """
    # Configure root logger for tests
    logging.basicConfig(
        level=logging.DEBUG,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    )

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


def cleanup_test_logging(module_names: list[str]) -> None:
    """Clean up logging configuration after tests.

    Args:
        module_names: List of module logger names to clean up.
    """
    for name in module_names:
        logging.getLogger(name).handlers.clear()


# List of module loggers to configure
MODULE_LOGGERS = [
    "nanodoc",  # Main logger
    "cli",  # V2 CLI
    "document",  # V2 Document
    "formatter",  # V2 Formatter
    "renderer",  # V2 Renderer
    "resolver",  # V2 Resolver
    "extractor",  # V2 Extractor
]
