#!/usr/bin/env node

var words = require('wordlist-english');
var list = words['english/10'].slice(0, 1000);

list.forEach(function(w) {
   console.log(w);
});
