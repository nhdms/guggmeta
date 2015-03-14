'use strict';

var app = angular.module('guggmeta', []);

app.controller('MainCtrl', function MainCtrl($scope, $timeout) {

  $scope.time = 'time...'
  var interval = 100;

  var tick = function () {
    $scope.time = Date.now();
    $timeout(tick, interval);
  }

  $timeout(tick, interval);

});
