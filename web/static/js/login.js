function login () {
	var email = $('#inputEmail').val();
	var pwd = $('#inputPassword').val();
	var rememberMe = $('#inputRemember').is(':checked');
	if (!email) { alert('Please enter the Email.'); return false; }
	if (!pwd) { alert('Please enter the Password.'); return false; }
	$.ajax({
	  	type: 'POST',
	  	dataType: 'json',
	  	url: 'api/user/login',
	  	data: {'email':email, 'pwd':pwd, 'remember':rememberMe ? 1:0},
	  	success: function (res, textStatus, jqXHR) {
			if (res.success) {
				var tourl = getURLParameterByName('returnurl');
				tourl = tourl || "/";
				window.location = tourl;
			} else {
				alert('login faild: ' + res.error);
			}
	  	}
	});
	return false;
}

function getURLParameterByName(name) {
    name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
    var regex = new RegExp("[\\?&]" + name + "=([^&#]*)"),
        results = regex.exec(location.search);
    return results === null ? "" : decodeURIComponent(results[1].replace(/\+/g, " "));
}

