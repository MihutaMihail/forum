const passwordInput = document.getElementById("password");
const confirmPasswordInput = document.getElementById("confirmPassword");

console.log("fdsfds");

document.getElementById("centerBoxForm").addEventListener("submit", function(event) {
    const password = passwordInput.value;
    const confirmPassword = confirmPasswordInput.value;

    console.log(password);
    console.log(confirmPassword);

    if (password !== confirmPassword) {
        event.preventDefault();
        setTimeout(function() {
            alert("Passwords do not match");
        }, 10);
    }
});