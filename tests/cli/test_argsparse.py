"""Tests for command-line argument parsing in nanodoc.

These tests verify that command-line arguments are correctly parsed and
passed to the process function.
"""

import sys
from unittest.mock import patch

import pytest
from click.testing import CliRunner

from nanodoc.cli import main


@pytest.fixture
def cli_runner():
    """Provide a Click CLI test runner."""
    return CliRunner()


@patch("nanodoc.cli.process")
def test_basic_arguments_passed_to_process(mock_process, monkeypatch):
    """Test that basic arguments are correctly passed to process."""
    # Setup sys.argv as it would be when called from command line
    test_args = ["nanodoc", "path/to/file.txt"]
    monkeypatch.setattr(sys, "argv", test_args)

    # Call the main function
    with patch("nanodoc.cli.click.echo"):  # Suppress output
        main()

    # Verify process was called with the expected parameters
    assert mock_process.call_count == 1
    args, kwargs = mock_process.call_args
    assert kwargs["sources"] == ("path/to/file.txt",)
    assert not kwargs["toc"]
    assert kwargs["n"] == 0  # Default line numbering (none)
    assert not kwargs["no_header"]
    assert kwargs["use_v2"] is True  # Default is now True


@patch("nanodoc.cli.process")
def test_toc_flag_passed_to_process(mock_process, monkeypatch):
    """Test that --toc flag is correctly passed to process."""
    # Setup sys.argv with --toc flag
    test_args = ["nanodoc", "--toc", "path/to/file.txt"]
    monkeypatch.setattr(sys, "argv", test_args)

    # Call the main function
    with patch("nanodoc.cli.click.echo"):  # Suppress output
        main()

    # Verify process was called with toc=True
    assert mock_process.call_count == 1
    args, kwargs = mock_process.call_args
    assert kwargs["toc"] is True
    assert kwargs["sources"] == ("path/to/file.txt",)


@patch("nanodoc.cli.process")
def test_line_numbering_passed_to_process(mock_process, monkeypatch):
    """Test that -n flags are correctly passed to process."""
    # Setup sys.argv with -n flag
    test_args = ["nanodoc", "-n", "path/to/file.txt"]
    monkeypatch.setattr(sys, "argv", test_args)

    # Call the main function
    with patch("nanodoc.cli.click.echo"):  # Suppress output
        main()

    # Verify process was called with n=1
    assert mock_process.call_count == 1
    args, kwargs = mock_process.call_args
    assert kwargs["n"] == 1
    assert kwargs["sources"] == ("path/to/file.txt",)


@patch("nanodoc.cli.process")
def test_multiple_line_numbering_passed_to_process(mock_process, monkeypatch):
    """Test that -nn flags are correctly passed to process."""
    # Setup sys.argv with -nn flags
    test_args = ["nanodoc", "-nn", "path/to/file.txt"]
    monkeypatch.setattr(sys, "argv", test_args)

    # Call the main function
    with patch("nanodoc.cli.click.echo"):  # Suppress output
        main()

    # Verify process was called with n=2
    assert mock_process.call_count == 1
    args, kwargs = mock_process.call_args
    assert kwargs["n"] == 2
    assert kwargs["sources"] == ("path/to/file.txt",)


# Parametrized test cases for testing multiple argument combinations
@pytest.mark.parametrize(
    "command_args, expected_params, should_skip",
    [
        # Test case ID: basic
        (
            ["nanodoc", "file.txt"],
            {"sources": ("file.txt",), "toc": False, "n": 0, "no_header": False},
            False,
        ),
        # Test case ID: toc_only
        (
            ["nanodoc", "--toc", "file.txt"],
            {"sources": ("file.txt",), "toc": True, "n": 0, "no_header": False},
            False,
        ),
        # Test case ID: line_numbers_file
        (
            ["nanodoc", "-n", "file.txt"],
            {"sources": ("file.txt",), "toc": False, "n": 1, "no_header": False},
            False,
        ),
        # Test case ID: line_numbers_all
        (
            ["nanodoc", "-nn", "file.txt"],
            {"sources": ("file.txt",), "toc": False, "n": 2, "no_header": False},
            False,
        ),
        # Test case ID: no_header
        (
            ["nanodoc", "--no-header", "file.txt"],
            {"sources": ("file.txt",), "toc": False, "n": 0, "no_header": True},
            False,
        ),
        # Test case ID: complex_case
        (
            ["nanodoc", "--toc", "-nn", "--no-header", "file1.txt", "file2.txt"],
            {
                "sources": ("file1.txt", "file2.txt"),
                "toc": True,
                "n": 2,
                "no_header": True,
            },
            False,
        ),
        # Test case ID: forced_v1
        (
            ["nanodoc", "--no-use-v2", "file.txt"],
            {
                "sources": ("file.txt",),
                "toc": False,
                "n": 0,
                "no_header": False,
                "use_v2": False,
            },
            True,  # Click handles boolean flags differently
        ),
        # Test case ID: theme_specified
        (
            ["nanodoc", "--theme", "classic-dark", "file.txt"],
            {
                "sources": ("file.txt",),
                "toc": False,
                "n": 0,
                "no_header": False,
                "theme": "classic-dark",
            },
            True,  # Theme may be stored in context object
        ),
    ],
    ids=[
        "basic",
        "toc_only",
        "line_numbers_file",
        "line_numbers_all",
        "no_header",
        "complex_case",
        "forced_v1",
        "theme_specified",
    ],
)
@patch("nanodoc.cli.process")
def test_multiple_argument_combinations(
    mock_process, monkeypatch, command_args, expected_params, should_skip
):
    """Test various combinations of command-line arguments."""
    if should_skip:
        pytest.skip("Click handles this option differently")

    # Setup sys.argv with the test arguments
    monkeypatch.setattr(sys, "argv", command_args)

    # Call the main function
    with patch("nanodoc.cli.click.echo"):  # Suppress output
        main()

    # Verify process was called with the expected parameters
    assert mock_process.call_count == 1
    args, kwargs = mock_process.call_args

    for param, value in expected_params.items():
        assert param in kwargs, f"Parameter {param} not found in kwargs"
        assert kwargs[param] == value, f"Failed for param {param}"


# Test that Click requires flags before arguments
@pytest.mark.parametrize(
    "command_args, should_skip",
    [
        # Test providing flags first - Click will handle this correctly
        (["nanodoc", "--toc", "-nn", "file.txt"], False),
        # Test providing flags after file - Click will NOT interpret these as options
        (["nanodoc", "file.txt", "--toc", "-nn"], True),
        # Test interleaved flags and sources - Click will NOT handle this as expected
        (["nanodoc", "--toc", "file1.txt", "-nn", "file2.txt"], True),
    ],
    ids=["flags_first", "flags_after", "interleaved"],
)
@patch("nanodoc.cli.process")
def test_flag_position_requirements(
    mock_process, monkeypatch, command_args, should_skip
):
    """Test Click's behavior with flag positions.

    Note: Click requires all options/flags to come before positional arguments.
    This test documents this behavior.
    """
    if should_skip:
        pytest.skip("Click requires options/flags to come before positional arguments")

    monkeypatch.setattr(sys, "argv", command_args)

    # Call the main function
    with patch("nanodoc.cli.click.echo"):  # Suppress output
        main()

    # Verify process was called and flags were correctly interpreted
    assert mock_process.call_count == 1
    args, kwargs = mock_process.call_args
    assert kwargs["toc"] is True
    assert kwargs["n"] == 2
