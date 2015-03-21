'use strict';

var app = angular.module('guggmeta', []);

app.controller('MainCtrl', function MainCtrl($scope, SubmissionService) {

  $scope.submissions = [];

  var load = function () {
    SubmissionService.getAll().then(function (response) {
      $scope.submissions = response.data.hits;
      console.log(response.data);
    });
  };

  load();

});

app.service('SubmissionService', function SubmissionService($http) {

  this.getOne = function (id) {
    return $http.get('/api/submission/' + id + '/');
  };

  this.getAll = function () {
    return $http.get('/api/submissions/');
  };

});
