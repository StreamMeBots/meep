var React = require('react'),
	helpers = require('./helpers');

module.exports = function(ctx) {
	React.render(
		<section className='content'>
			sup yo!
		</section>,
		helpers.getPageDiv()
	);
}
