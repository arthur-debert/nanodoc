"""Boot configuration module for nanodoc v2.

This module handles core configuration like logging setup.
"""

import logging
import os
import sys


def configure_logging(verbose: bool = False) -> None:
    """Configure logging for nanodoc v2.

    Logging can be enabled in two ways:
    1. Setting verbose=True when calling this function
    2. Setting the NANODOC_VERBOSE environment variable

    Args:
        verbose: Whether to enable verbose (DEBUG) logging.
    """
    # Check both the verbose flag and environment variable
    is_verbose = verbose or os.environ.get("NANODOC_VERBOSE") is not None

    # Set log level based on verbose flag
    log_level = logging.DEBUG if is_verbose else logging.WARNING

    # Configure root logger
    logging.basicConfig(
        level=log_level,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    )

    # Configure all module loggers
    for name in MODULE_LOGGERS:
        module_logger = logging.getLogger(name)
        module_logger.setLevel(log_level)
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
