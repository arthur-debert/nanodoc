from pathlib import Path

from nanodoc.v1.core import run_all
from nanodoc.v1.data import ContentItem, LineRange
from nanodoc.v1.files import create_content_item
from nanodoc.v1.toc import (
    calculate_line_numbers,
    calculate_toc_size,
    create_toc_content,
    format_filenames,
    generate_table_of_contents,
    group_content_items_by_file,
)


def test_group_content_items_by_file(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")

    content_items = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    file_groups = group_content_items_by_file(content_items)

    # Check that we have the correct number of groups
    assert len(file_groups) == 2

    # Check that each file path is a key in the dictionary
    assert str(test_file1) in file_groups
    assert str(test_file2) in file_groups

    # Check that each group contains the correct ContentItem
    assert len(file_groups[str(test_file1)]) == 1
    assert len(file_groups[str(test_file2)]) == 1
    assert file_groups[str(test_file1)][0] == content_items[0]
    assert file_groups[str(test_file2)][0] == content_items[1]


def test_calculate_toc_size(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")

    content_items = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    file_groups = group_content_items_by_file(content_items)
    toc_size = calculate_toc_size(file_groups)

    # TOC size should be:
    # 2 lines for header (header + blank line)
    # 2 lines for entries (1 per file)
    # 1 line for footer (blank line)
    # Total: 5 lines
    assert toc_size == 5

    # Test with multiple ranges for the same file
    test_file3 = tmpdir.join("test_file3.txt")
    test_file3.write("Line 5\nLine 6\nLine 7\nLine 8")

    # Create two ContentItems for the same file
    content_items = [
        create_content_item(str(test_file3)),
        create_content_item(str(test_file3)),
    ]

    file_groups = group_content_items_by_file(content_items)
    toc_size = calculate_toc_size(file_groups)

    # TOC size should be:
    # 2 lines for header (header + blank line)
    # 1 line for the file entry
    # 2 lines for subentries (1 per range)
    # 1 line for footer (blank line)
    # Total: 6 lines
    assert toc_size == 6


def test_calculate_line_numbers(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")

    content_items = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    file_groups = group_content_items_by_file(content_items)
    toc_size = calculate_toc_size(file_groups)
    line_numbers = calculate_line_numbers(file_groups, toc_size)

    # Check that we have line numbers for both files
    assert str(test_file1) in line_numbers
    assert str(test_file2) in line_numbers

    # First file should start at line toc_size + 3
    assert line_numbers[str(test_file1)] == toc_size + 3

    # Second file should start after first file
    # First file has 2 lines of content + 3 lines for header/footer
    # So second file should start at toc_size + 3 + 2 + 3
    assert line_numbers[str(test_file2)] == toc_size + 3 + 2 + 3


def test_format_filenames(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")

    content_items = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    file_groups = group_content_items_by_file(content_items)

    # Test with default style (None)
    formatted_filenames = format_filenames(file_groups)
    assert formatted_filenames[str(test_file1)] == "test_file1.txt"
    assert formatted_filenames[str(test_file2)] == "test_file2.txt"

    # Test with 'nice' style
    formatted_filenames = format_filenames(file_groups, style="nice")
    assert formatted_filenames[str(test_file1)] == ("Test File1 (test_file1.txt)")
    assert formatted_filenames[str(test_file2)] == ("Test File2 (test_file2.txt)")

    # Test with 'filename' style
    formatted_filenames = format_filenames(file_groups, style="filename")
    assert formatted_filenames[str(test_file1)] == "test_file1.txt"
    assert formatted_filenames[str(test_file2)] == "test_file2.txt"

    # Test with 'path' style
    formatted_filenames = format_filenames(file_groups, style="path")
    assert formatted_filenames[str(test_file1)] == str(test_file1)
    assert formatted_filenames[str(test_file2)] == str(test_file2)


def test_create_toc_content(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")

    content_items = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    file_groups = group_content_items_by_file(content_items)
    toc_size = calculate_toc_size(file_groups)
    line_numbers = calculate_line_numbers(file_groups, toc_size)
    formatted_filenames = format_filenames(file_groups)

    toc = create_toc_content(file_groups, formatted_filenames, line_numbers)

    # Check that the TOC contains the expected content
    assert "TOC" in toc
    assert "test_file1.txt" in toc
    assert "test_file2.txt" in toc
    assert str(line_numbers[str(test_file1)]) in toc
    assert str(line_numbers[str(test_file2)]) in toc

    # Test with multiple ranges for the same file
    test_file3 = tmpdir.join("test_file3.txt")
    test_file3.write("Line 5\nLine 6\nLine 7\nLine 8")

    # Create two ContentItems for the same file
    content_items = [
        create_content_item(str(test_file3)),
        create_content_item(str(test_file3)),
    ]

    file_groups = group_content_items_by_file(content_items)
    toc_size = calculate_toc_size(file_groups)
    line_numbers = calculate_line_numbers(file_groups, toc_size)
    formatted_filenames = format_filenames(file_groups)

    toc = create_toc_content(file_groups, formatted_filenames, line_numbers)

    # Check that the TOC contains the expected content
    assert "TOC" in toc
    assert "test_file3.txt" in toc
    # Since we're not using line ranges, we won't have a. and b. subentries
    # Instead, check that the file name appears and the line number is present
    assert str(line_numbers[str(test_file3)]) in toc


def test_generate_table_of_contents(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")
    file_paths = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    toc, line_numbers = generate_table_of_contents(file_paths)

    # Check TOC content
    assert "TOC" in toc
    assert "test_file1.txt" in toc
    assert "test_file2.txt" in toc

    # Check line numbers
    assert isinstance(line_numbers, dict)
    file1_path = str(test_file1)
    file2_path = str(test_file2)
    assert line_numbers[file1_path] > 0
    assert line_numbers[file2_path] > line_numbers[file1_path]


def test_generate_table_of_contents_with_style(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")
    file_paths = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    # Test with 'nice' style
    toc, line_numbers = generate_table_of_contents(file_paths, style="nice")

    # Check that styled filenames are in the TOC
    assert "Test File1 (test_file1.txt)" in toc
    assert "Test File2 (test_file2.txt)" in toc

    # Test with 'filename' style
    toc, _ = generate_table_of_contents(file_paths, style="filename")

    # Check that plain filenames are in the TOC
    assert "test_file1.txt" in toc
    assert "test_file2.txt" in toc

    # Test with 'path' style
    toc, _ = generate_table_of_contents(file_paths, style="path")

    # Check that full paths are in the TOC
    assert str(test_file1) in toc
    assert str(test_file2) in toc


# The following tests use process_all but are kept for integration testing
# They verify that the TOC generation works correctly when integrated with
# the rest of the code


def test_process_all_toc_generation():
    """Test TOC generation with multiple content items."""
    content_items = [
        ContentItem(
            original_arg="test1.txt",
            file_path="test1.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 1\n", "Line 2\n"],
        ),
        ContentItem(
            original_arg="test2.txt",
            file_path="test2.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 3\n", "Line 4\n"],
        ),
    ]

    result = run_all(
        content_items,
        line_number_mode=None,
        generate_toc=True,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
    )

    assert "Table of Contents" in result
    assert "test1.txt" in result
    assert "test2.txt" in result
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result


def test_process_all_with_no_header():
    """Test processing multiple content items without headers."""
    content_items = [
        ContentItem(
            original_arg="test1.txt",
            file_path="test1.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 1\n", "Line 2\n"],
        ),
        ContentItem(
            original_arg="test2.txt",
            file_path="test2.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 3\n", "Line 4\n"],
        ),
    ]

    result = run_all(
        content_items,
        line_number_mode=None,
        generate_toc=False,
        theme=None,
        show_header=False,
        sequence=None,
        style=None,
    )

    assert "test1.txt" not in result
    assert "test2.txt" not in result
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result


def test_process_all_with_header_sequence():
    """Test processing multiple content items with header sequence."""
    content_items = [
        ContentItem(
            original_arg="test1.txt",
            file_path="test1.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 1\n", "Line 2\n"],
        ),
        ContentItem(
            original_arg="test2.txt",
            file_path="test2.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 3\n", "Line 4\n"],
        ),
    ]

    result = run_all(
        content_items,
        line_number_mode=None,
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence="numerical",
        style=None,
    )

    assert "1. test1.txt" in result
    assert "2. test2.txt" in result
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result


def test_process_all_with_header_style():
    """Test processing multiple content items with header style."""
    content_items = [
        ContentItem(
            original_arg="test1.txt",
            file_path="test1.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 1\n", "Line 2\n"],
        ),
        ContentItem(
            original_arg="test2.txt",
            file_path="test2.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 3\n", "Line 4\n"],
        ),
    ]

    result = run_all(
        content_items,
        line_number_mode=None,
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence=None,
        style="nice",
    )

    assert "Test1 (test1.txt)" in result
    assert "Test2 (test2.txt)" in result
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result


def test_group_by_file(fixture_content_item):
    """Test grouping content items by file."""
    content_items = [fixture_content_item]
    file_groups = group_content_items_by_file(content_items)

    assert fixture_content_item.file_path in file_groups
    assert len(file_groups[fixture_content_item.file_path]) == 1
    assert file_groups[fixture_content_item.file_path][0] == fixture_content_item


def test_calculate_line_numbers(fixture_content_item):
    """Test calculating line numbers for files."""
    content_items = [fixture_content_item]
    file_groups = group_content_items_by_file(content_items)
    toc_size = 5

    line_numbers = calculate_line_numbers(file_groups, toc_size)

    assert fixture_content_item.file_path in line_numbers
    assert line_numbers[fixture_content_item.file_path] == toc_size + 3


def test_format_filenames(fixture_content_item):
    """Test formatting filenames with different styles."""
    content_items = [fixture_content_item]
    file_groups = group_content_items_by_file(content_items)

    # Test default style
    formatted = format_filenames(file_groups, None)
    assert (
        formatted[fixture_content_item.file_path] == fixture_content_item.original_arg
    )

    # Test 'nice' style
    formatted = format_filenames(file_groups, "nice")
    expected = f"Test {Path(fixture_content_item.original_arg).stem.title()} ({fixture_content_item.original_arg})"
    assert formatted[fixture_content_item.file_path] == expected

    # Test 'simple' style
    formatted = format_filenames(file_groups, "simple")
    assert (
        formatted[fixture_content_item.file_path] == fixture_content_item.original_arg
    )

    # Test 'full' style
    formatted = format_filenames(file_groups, "full")
    assert formatted[fixture_content_item.file_path] == fixture_content_item.file_path


def test_generate_table_of_contents(fixture_content_item):
    """Test generating table of contents."""
    content_items = [fixture_content_item]
    toc, line_numbers = generate_table_of_contents(content_items, None)

    assert fixture_content_item.original_arg in toc
    assert str(line_numbers[fixture_content_item.file_path]) in toc
