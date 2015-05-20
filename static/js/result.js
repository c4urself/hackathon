$(function() {
	var photos = $(".photo"),
	    mosaic = $(".mosaic .mat img"),
	    search = $(".search");

	function switchPhoto(photo, speed) {
		mosaic.fadeOut(speed, function() {
			mosaic.attr("src", photo.data("mosaic"));
			mosaic.fadeIn(speed);
		});
		photos.removeClass("active");
		photo.addClass("active");
	}

	photos.click(function() {
		switchPhoto($(this), 300);
	});

	search.submit(function(e) {
		e.preventDefault();
		document.location = "/mosaic/" + search.find("input").val();
	})

	switchPhoto(photos.first(), 0);
});
