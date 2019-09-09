$('div[class^="grid-symbol"]').each(function (index) { 
	if ((index+1) % 10 == 0) {
		$(this).css("border-right", "3px solid black");
	} else if ((index+1) % 5 == 0) {
		$(this).css("border-right", "2px solid black");
	} 

	// bold border for every 5th and 10th row
	if (index % 500 > 450) {
		$(this).css("border-bottom", "3px solid black");
	} else if (index % 250 >= 200) {
		$(this).css("border-bottom", "2px solid black");
	}
});
