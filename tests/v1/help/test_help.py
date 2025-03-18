import os
import subprocess
import sys

from nanodoc.v1.help import get_options_section, get_topics_section

# Get the parent directory of the current module
MODULE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# Use Python module approach instead of direct script execution
PYTHON_CMD = sys.executable
NANODOC_MODULE = "src.nanodoc"


def test_help():
    result = subprocess.run(
        [PYTHON_CMD, "-m", NANODOC_MODULE, "help"], capture_output=True, text=True
    )
    assert result.returncode == 0
    assert "# nanodoc" in result.stdout


def test_no_args():
    result = subprocess.run(
        [PYTHON_CMD, "-m", NANODOC_MODULE], capture_output=True, text=True
    )
    assert result.returncode == 0
    assert "Usage: python -m src.nanodoc" in result.stdout
    assert "# nanodoc" not in result.stdout


def test_get_options_section():
    """Test getting the options section."""
    options = get_options_section()
    assert isinstance(options, str)
    assert "OPTIONS:" in options
    assert "--toc" in options
    assert "--no-header" in options
    assert "--sequence" in options
    assert "--style" in options
    assert "-n" in options
    assert "--txt-ext" in options


def test_get_topics_section():
    """Test getting the topics section."""
    topics = get_topics_section()
    assert isinstance(topics, str)
    assert "TOPICS:" in topics
    assert "line numbers" in topics.lower()
    assert "table of contents" in topics.lower()
    assert "file headers" in topics.lower()
