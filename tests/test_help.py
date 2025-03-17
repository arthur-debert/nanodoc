import os
import re
import subprocess
import sys

from nanodoc.help import get_options_section, get_topics_section

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
    assert "usage: nanodoc" in result.stdout
    assert "# nanodoc" not in result.stdout


def test_get_options_section():
    """Test that get_options_section returns the formatted options."""
    options_content = get_options_section()

    # Check that the options section contains expected options
    assert "-v:" in options_content
    assert "--toc:" in options_content
    assert "--no-header:" in options_content
    assert "--sequence:" in options_content
    assert "--style:" in options_content

    # Check the formatting
    # Each option should be followed by spaces and then the help text
    assert re.search(r"-v:\s+Enable verbose mode", options_content)
    assert re.search(r"--toc:\s+Generate table of contents", options_content)


def test_get_topics_section():
    """Test that get_topics_section returns the formatted topics."""
    topics_content = get_topics_section()

    # Check that the topics section contains expected guides
    assert "manifesto:" in topics_content
    assert "quickstart:" in topics_content

    # Check the formatting
    # Each topic should be followed by spaces and then the description
    assert re.search(r"manifesto:\s+Less clutter, less distraction", topics_content)
    assert "quickstart:" in topics_content
