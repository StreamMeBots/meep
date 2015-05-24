var React = require('react'),
	cookie = require('cookie'),
	xhr = require('xhr'),
	user = require('../user');

module.exports = function(el, cb) {
	var sessid = cookie.parse(document.cookie).sessid;

	if(!sessid) {
		return cb();
	}

	xhr({
		method: 'GET',
		uri: '/api/me',
		json: true
	}, function(err, resp, body) {
		if(err) {
			return cb();
		}

		if(resp.statusCode !== 200) {
			return cb();
		}

		user.set(body);
		cb();
	});
}
