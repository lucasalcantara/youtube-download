$(document).ready(function(){
	$("#submit").click(function(){
		if(validateUrlValue() && validateFormat()){
			$("#loader").show();

			$.ajax({
			  type: "POST",
			  url: "post",
			  data: {
			  	urls: $('input[name=url]').serialize(), 
			  	format: $("input[name=format]:checked").val(),
			  	upload: $("input[name=upload]").is(':checked'), 
			  	option: $(".nav.nav-tabs li[class=active] a ").text()
			  },
			  success: function(data){
			  	alert(data);

			  	$("#loader").hide();
			  }
			});
		}
	});

	$(".panel-body li").click(function(){
	     $("#url").val("");
	     $('input[name=url]').not("#url").parent().remove();
	});

	$(".glyphicon-plus-sign").click(function(){
		var field = '<div>';
		field += '	URL: <input type="text" name="url" style="width: 80%;" /> '
		field += '	<span class="glyphicon glyphicon-minus-sign" aria-hidden="true"></span>';
		field += "<br></div>";

		$("#add").append(field);
	});

	$("#add").on('click', ".glyphicon-minus-sign", function(){
		$(this).parent().remove();
	});
});

function validateUrlValue(){
	var validated = true;
	var values = [];
	var expression = /[-a-zA-Z0-9@:%_\+.~#?&//=]{2,256}\.[a-z]{2,4}\b(\/[-a-zA-Z0-9@:%_\+.~#?&//=]*)?/gi;
	var regex = new RegExp(expression);

	$('input[name=url]').each(function(){
		var url = $(this).val();
		values.push(url);

		if (!url.match(regex)) {
		    alert("Url value is not a youtube url.");
		    validated = false;
		}
	});

	return validated && !existsDuplicateValues(values);
}

function existsDuplicateValues(arr){
	var sorted_arr = arr.slice().sort();
	var results = [];
	for (var i = 0; i < arr.length - 1; i++) {
	    if (sorted_arr[i + 1] == sorted_arr[i]) {
	        results.push(sorted_arr[i]);
	    }
	}

	if(results.length){
		alert("Duplicate url detected.")
	}

	return results.length;
}

function validateFormat(){
	if(!$("input[name=format]:checked").length){
		alert("Please select a format");
		return false;
	}

	return true;
}