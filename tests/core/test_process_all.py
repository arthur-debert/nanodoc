from nanodoc.core import process_all
from nanodoc.files import create_content_item


def test_process_all_single_file(tmpdir):
    # Test processing a single file
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")
    file_path = str(test_file)

    content_item = create_content_item(file_path)
    result = process_all([content_item])

    # Check that the result contains the file header and content
    assert "test_file.txt" in result
    assert "Line 1" in result
    assert "Line 5" in result


def test_process_all_multiple_files(tmpdir):
    # Test processing multiple files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("File 1 Line 1\nFile 1 Line 2\nFile 1 Line 3")
    file_path1 = str(test_file1)

    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("File 2 Line A\nFile 2 Line B\nFile 2 Line C")
    file_path2 = str(test_file2)

    content_items = [create_content_item(file_path1), create_content_item(file_path2)]
    result = process_all(content_items)

    # Check that the result contains both file headers and content
    assert "test_file1.txt" in result
    assert "File 1 Line 1" in result
    assert "test_file2.txt" in result
    assert "File 2 Line A" in result


def test_process_all_with_line_numbers(tmpdir):
    # Test processing with line numbers
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")
    file_path = str(test_file)

    content_item = create_content_item(file_path)

    # Test with file line numbers
    result = process_all([content_item], line_number_mode="file")
    assert "   1: Line 1" in result
    assert "   5: Line 5" in result

    # Test with global line numbers
    result = process_all([content_item], line_number_mode="all")
    assert "   1: Line 1" in result
    assert "   5: Line 5" in result


def test_process_all_with_toc(tmpdir):
    # Test processing with table of contents
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("File 1 Line 1\nFile 1 Line 2\nFile 1 Line 3")
    file_path1 = str(test_file1)

    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("File 2 Line A\nFile 2 Line B\nFile 2 Line C")
    file_path2 = str(test_file2)

    content_items = [create_content_item(file_path1), create_content_item(file_path2)]
    result = process_all(content_items, generate_toc=True)

    # Check that the result contains the TOC
    assert "TOC" in result
    assert "test_file1.txt" in result
    assert "test_file2.txt" in result


def test_process_all_with_partial_content(tmpdir):
    # Test processing with partial content
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")
    file_path = str(test_file)

    # Create a ContentItem with a specific line range
    content_item = create_content_item(f"{file_path}:L2-4")
    result = process_all([content_item])

    # Check that the result contains only the specified lines
    assert "test_file.txt" in result
    assert "Line 1" not in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result
    assert "Line 5" not in result


def test_process_all_with_multiple_ranges(tmpdir):
    # Test processing with multiple ranges from the same file
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")
    file_path = str(test_file)

    # Create ContentItems with different ranges from the same file
    content_items = [
        create_content_item(f"{file_path}:L1-2"),
        create_content_item(f"{file_path}:L4-5"),
    ]
    result = process_all(content_items)

    # Check that the result contains both ranges
    assert "test_file.txt" in result
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" not in result
    assert "Line 4" in result
    assert "Line 5" in result


def test_process_all_with_toc_and_multiple_ranges(tmpdir):
    # Test processing with TOC and multiple ranges
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")
    file_path = str(test_file)

    # Create ContentItems with different ranges from the same file
    content_items = [
        create_content_item(f"{file_path}:L1-2"),
        create_content_item(f"{file_path}:L4-5"),
    ]
    result = process_all(content_items, generate_toc=True)

    # Check that the result contains the TOC with subentries
    assert "TOC" in result
    assert "test_file.txt" in result
    assert "a. L1-2" in result
    assert "b. L4-5" in result
