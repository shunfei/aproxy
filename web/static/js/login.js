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

$(document).ready(function(){
	$.get("api/oauth/list",function(data, status){
		console.dir("Data: " + data + "\nStatus: " + status);
		if (data.data && data.data.length > 0) {
			$("#oauthLogin").show();
			var tourl = getURLParameterByName('returnurl');
			var oauthProviders = "";
			for (var i = 0; i < data.data.length; i++) {
				var provider = data.data[i];
				oauthProviders += '<a href="api/oauth/login?provider='+provider+'&returnurl='+tourl+'" class="btn btn-info">'+provider+'</a>';
			}
			$("#oauthProviderList").html(oauthProviders)
		}
	});
});

