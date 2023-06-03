//
// Confirm Password
//

const passwordInput = document.getElementById("password");
const confirmPasswordInput = document.getElementById("confirmPassword");

document.getElementById("centerBoxForm").addEventListener("submit", function(event) {
    const password = passwordInput.value;
    const confirmPassword = confirmPasswordInput.value;

    if (password !== confirmPassword) {
        event.preventDefault();
        setTimeout(function() {
            alert("Passwords do not match");
        }, 10);
    }
});

//
// Password Validation
//

var form = document.getElementById("centerBoxForm");

form.addEventListener("submit", function (event) {
  event.preventDefault();

  // Get the values entered by the user
  var email = document.getElementById("mail").value;
  var username = document.getElementById("username").value;
  var password = document.getElementById("password").value;

  // Define regular expressions for validation
  var emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  var usernameRegex = /^.{5,}$/;
  var passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/;

  // Validate the email
  if (!emailRegex.test(email)) {
    alert("Please enter a valid email address");
    return;
  }

  // Validate the username length
  if (!usernameRegex.test(username)) {
    alert("Username must be at least 5 characters long");
    return;
  }

  // Validate the password requirements
  if (!passwordRegex.test(password)) {
    alert("Password must contain at least 8 characters, including a capital letter, a number, and a special character");
    return;
  }

  // Submit form if all tests are validated
  form.submit();
});