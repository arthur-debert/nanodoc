import os
import subprocess
import sys

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
