
$('div[class^="grid-symbol"]').each(function (index) { 
	if ((index+1) % 10 == 0) {
		$(this).css("border-right", "3px solid black");
	} else if ((index+1) % 5 == 0) {
		$(this).css("border-right", "2px solid black");
	} 

	// bold border for every 5th and 10th row
	if (index % 500 >= 450) {
		$(this).css("border-bottom", "3px solid black");
	} else if (index % 250 >= 200) {
		$(this).css("border-bottom", "2px solid black");
	}
});

function writePDF() {
	//var source = $('div[class^="grid-container"]')[0];
	var source = $('body,html')[0];
	var doc = new jsPDF('p', 'in', 'letter');
	doc.fromHTML(
		source, // HTML string or DOM elem ref.
		0.5, // x coord
		0.5, // y coord
		{	
			'width': 7.5 // max width of content on PDF
	});

	doc.save('test.pdf');
}
