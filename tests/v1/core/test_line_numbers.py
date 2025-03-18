import os

from nanodoc.v1.core import run_file
from nanodoc.v1.data import ContentItem, LineRange
from nanodoc.v1.files import create_content_item


def test_process_file_no_line_numbers(tmpdir):
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2")
    file_path = str(test_file)
    output, _ = run_file(create_content_item(file_path), None, 0)
    assert "Line 1" in output
    assert "Line 2" in output
    assert "1:" not in output
    assert "2:" not in output


def test_process_file_with_line_numbers_all(tmpdir):
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2")
    file_path = str(test_file)
    output, _ = run_file(create_content_item(file_path), "all", 0)
    assert "1: Line 1" in output
    assert "2: Line 2" in output


def test_process_file_with_line_numbers_file(tmpdir):
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2")
    file_path = str(test_file)
    output, _ = run_file(create_content_item(file_path), "file", 0)
    assert "1: Line 1" in output
    assert "2: Line 2" in output


def test_process_file_not_found():
    file_path = "nonexistent_file.txt"
    output, _ = run_file(create_content_item(file_path), None, 0)
    assert "Error: File not found" in output


def test_process_file_header_assignment(tmpdir):
    test_file = tmpdir.join("test_file.txt")
    test_file.write("test")
    file_path = str(test_file)
    output, _ = run_file(create_content_item(file_path), None, 0)
    header = "\n" + os.path.basename(file_path) + "\n\n"
    assert output.startswith(header)


def test_process_file_with_no_header(tmpdir):
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2")
    file_path = str(test_file)
    output, _ = run_file(create_content_item(file_path), None, 0, show_header=False)
    assert "Line 1" in output
    assert "Line 2" in output
    assert not output.startswith("\n")
    assert os.path.basename(file_path) not in output.split("\n")[0]


def test_process_file_with_header_sequence(tmpdir):
    test_file = tmpdir.join("test_file.txt")
    test_file.write("test")
    file_path = str(test_file)

    # Test numerical sequence
    output, _ = run_file(
        create_content_item(file_path), None, 0, sequence="numerical", seq_index=0
    )
    assert "\n1. test_file.txt\n\n" in output

    # Test letter sequence
    output, _ = run_file(
        create_content_item(file_path), None, 0, sequence="letter", seq_index=1
    )
    assert "\nb. test_file.txt\n\n" in output

    # Test roman sequence
    output, _ = run_file(
        create_content_item(file_path), None, 0, sequence="roman", seq_index=2
    )
    assert "\niii. test_file.txt\n\n" in output


def test_process_file_with_header_style(tmpdir):
    test_file = tmpdir.join("test_file.txt")
    test_file.write("test")
    file_path = str(test_file)

    # Test nice style with sequence
    output, _ = run_file(
        create_content_item(file_path),
        None,
        0,
        sequence="numerical",
        seq_index=0,
        style="nice",
    )
    assert "\n1. Test File (test_file.txt)\n\n" in output


def test_line_numbers_file_mode():
    """Test line numbering in file mode."""
    content_item = ContentItem(
        file_path="test.txt",
        ranges=[LineRange(start=1, end=3)],
        content=["Line 1\n", "Line 2\n", "Line 3\n"],
    )

    result, _ = run_file(content_item, "file", 0)

    assert "1: Line 1" in result
    assert "2: Line 2" in result
    assert "3: Line 3" in result


def test_line_numbers_all_mode():
    """Test line numbering in all mode."""
    content_item = ContentItem(
        file_path="test.txt",
        ranges=[LineRange(start=1, end=3)],
        content=["Line 1\n", "Line 2\n", "Line 3\n"],
    )

    result, _ = run_file(content_item, "all", 10)

    assert "11: Line 1" in result
    assert "12: Line 2" in result
    assert "13: Line 3" in result


def test_line_numbers_none_mode():
    """Test line numbering in none mode."""
    content_item = ContentItem(
        file_path="test.txt",
        ranges=[LineRange(start=1, end=3)],
        content=["Line 1\n", "Line 2\n", "Line 3\n"],
    )

    result, _ = run_file(content_item, None, 0)

    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "1:" not in result
    assert "2:" not in result
    assert "3:" not in result
