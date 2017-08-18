$(document).ready(function(){
	$("#submit").click(function(){
		if(validateUrlValue() && validateFormat()){
			$("#loader").show();

			$.ajax({
			  type: "POST",
			  url: "post",
			  data: {
			  	url: $("#url").val(), 
			  	format: $("input[name=format]:checked").val(), 
			  	option: $(".nav.nav-tabs li[class=active] a ").text()
			  },
			  success: function(data){
			  	alert(data);

			  	$("#loader").hide();
			  }
			});
		}
	});
});

function validateUrlValue(){
	var expression = /[-a-zA-Z0-9@:%_\+.~#?&//=]{2,256}\.[a-z]{2,4}\b(\/[-a-zA-Z0-9@:%_\+.~#?&//=]*)?/gi;
	var regex = new RegExp(expression);
	var url = $("#url").val();

	if (!url.match(regex)) {
	    alert("Url value is not a youtube url.");
	    return false;
	}

	return true;
}

function validateFormat(){
	if(!$("input[name=format]:checked").length){
		alert("Please select a format");
		return false;
	}

	return true;
}