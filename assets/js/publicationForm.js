//
// Image
//

var imageInputVisible = false;

// Toggles file input
function toggleFileInput() {
    const fileInput = document.getElementById('fileInput');
    var addImageBoolean = document.getElementById('addImageBoolean');

    if (imageInputVisible) {
        fileInput.style.display = 'none';
        document.getElementById('toggleFile').textContent = "Add Image";
        addImageBoolean.value = "false";
    } else {
        fileInput.style.display = 'block';
        document.getElementById('toggleFile').textContent = "Remove Image";
        addImageBoolean.value = "true";
    }
    imageInputVisible = !imageInputVisible;
}

// Hides default broken image icon
function hideBrokenImageIcon() {
    document.getElementById('imagePreview').style.display = 'none';
}

// Display current image chosen by the user
function displayImage() {
    const file = document.getElementById('image').files[0];
    const imagePreview = document.getElementById('imagePreview');

    if (!file.type.startsWith('image/')) {
        alert('Please select an image file');
        document.getElementById('image').value = '';
    } else {
        const reader = new FileReader();
        reader.onload = function() {
            imagePreview.src = reader.result;
        }
        reader.readAsDataURL(file);
        imagePreview.style.display = 'block';
    }
}

//
// Tags
//

// Add new tag
function addTag(selectedTag) {
    var selectedTagsDiv = document.getElementById('selectedTags');
    var selectedTagsList = document.getElementById('selected-tags');
    var selectedTagsInput = document.getElementById('selected-tags-input');

    if (selectedTagsList.childElementCount  === 0) {
        selectedTagsDiv.style.display = 'block'; 
    }

    var tagExists = Array.from(
        selectedTagsList.getElementsByTagName("li")).some(li => li.innerText === selectedTag || 
        (li.getElementsByTagName("button")[0] && li.getElementsByTagName("button")[0].innerText === "x" && 
        li.innerText === selectedTag + "x")
    );
        
    if (!tagExists) {
        var newTag = document.createElement("li");
        newTag.innerText = selectedTag;
        selectedTagsList.appendChild(newTag);

        // Add delete button
        var deleteButton = document.createElement("button");
        deleteButton.innerText = "x";
        deleteButton.classList.add("deleteTags-button");
        deleteButton.onclick = function() {
            newTag.remove();
            var selectedTags = Array.from(selectedTagsList.getElementsByTagName("li")).map(li => li.innerText);
            selectedTagsInput.value = JSON.stringify(selectedTags);

            if (selectedTags.length === 0) {
                selectedTagsDiv.style.display = 'none'; 
            }
        };
        newTag.appendChild(deleteButton);

        // Add tag to list
        var selectedTags = Array.from(selectedTagsList.getElementsByTagName("li")).map(li => li.innerText);
        selectedTagsInput.value = JSON.stringify(selectedTags);
    }
}

// After the html document has loaded, execute script
// The script is getting the list of tags to then add it to the list
document.addEventListener('DOMContentLoaded', function() {
    var existingTags = document.getElementById('existing-tags'); 
    var existingTagsValue = existingTags.value;
    
    // Remove the first and last characters if they are '[' and ']'
    if (existingTagsValue.charAt(0) === '[' && existingTagsValue.charAt(existingTagsValue.length - 1) === ']') {
        existingTagsValue = existingTagsValue.substring(1, existingTagsValue.length - 1);
    }
    var tagsArray = existingTagsValue.split(" ");

    if (tagsArray != "") {
        tagsArray.forEach(function(tag) {
            addTag(tag);
        });
    }
});

// Get the current selected option in our tags list
function getSelectedValue() {
    return document.getElementById('tags').value;
}

//
// HELP
//

// Array.from(selectedTagsList.getElementsByTagName("li")) --- from all these elements (li)
// .map(li => li.innerText)                                --- get me the value of li elements (innerText)
// .includes(selectedTag)                                  --- check if this value is present in the array
