import re
import subprocess
import sys

import nanodoc

# Use Python module approach for running nanodoc
PYTHON_CMD = sys.executable
NANODOC_MODULE = "src.nanodoc"


def test_version_import():
    """Test that the version can be imported from the package."""
    assert hasattr(nanodoc, "__version__")
    assert isinstance(nanodoc.__version__, str)
    assert re.match(r"^\d+\.\d+\.\d+$", nanodoc.__version__)


def test_version_flag():
    """Test that the --version flag works correctly."""
    result = subprocess.run(
        [PYTHON_CMD, "-m", NANODOC_MODULE, "--version"], capture_output=True, text=True
    )

    # The version flag exits with code 0
    assert result.returncode == 0

    # Check that the version output matches the expected format
    # The output could be in either stdout or stderr
    # depending on how argparse is configured
    version_output = result.stdout.strip() or result.stderr.strip()
    assert "nanodoc" in version_output
    assert nanodoc.__version__ in version_output


def test_version_module():
    """Test that the version module works correctly."""
    from nanodoc.version import VERSION, get_version

    assert VERSION == nanodoc.__version__
    assert get_version() == VERSION
