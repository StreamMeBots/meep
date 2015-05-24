var React = require('react'),
	helpers = require('./helpers'),
	Landing = require('../components/landing.jsx');

module.exports = function(ctx) {
	React.render(
		<Landing />,
		helpers.getPageDiv()
	);
}
