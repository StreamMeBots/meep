var User = function() {
	this.attributes = {
		authenticated: false
	};
};

User.prototype.set = function(attributes) {
	this.attributes = attributes;
	this.attributes.authenticated = true;
};

User.prototype.isAuthenticated = function() {
	return this.attributes.authenticated;
}

User.prototype.get = function(k) {
	switch(k) {
		case 'avatar':
			return this.attributes._links.avatar.href || '';
		default:
			return this.attributes[k] || '';
	}
}

module.exports = new User();
