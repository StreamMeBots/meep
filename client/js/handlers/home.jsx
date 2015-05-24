var React = require('react'),
	helpers = require('./helpers'),
	Home = require('../components/home.jsx');

module.exports = function(ctx) {
	React.render(
		<Home />,
		helpers.getPageDiv()
	);
}
