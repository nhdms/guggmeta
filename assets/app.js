'use strict';

var app = angular.module('guggmeta', ['ngMaterial', 'ngRoute', 'ngAria', 'ngAnimate']);

app.config(['$routeProvider', '$locationProvider', function ($routeProvider, $locationProvider) {
  $routeProvider
    .when('/submissions', {
      templateUrl: '/assets/partials/submission-list.tmpl.html',
      controller: 'SubmissionListCtrl'
    })
    .when('/submissions/:id', {
      templateUrl: '/assets/partials/submission-detail.tmpl.html',
      controller: 'SubmissionDetailCtrl'
    })
    .otherwise('/submissions');
  $locationProvider.html5Mode({
    enabled: true,
    requireBase: false
  });
}]);

app.controller('SubmissionListCtrl', ['$scope', 'SubmissionService', function ($scope, SubmissionService) {
  SubmissionService.getAll().then(function (response) {
    $scope.submissions = response.data.results;
  });
}]);

app.controller('SubmissionDetailCtrl', ['$scope', '$routeParams', 'SubmissionService', function ($scope, $routeParams, SubmissionService) {
  $scope.id = $routeParams.id;
  SubmissionService.getOne($routeParams.id).then(function (response) {
    $scope.submission = response.data;
  });
}]);

app.service('SubmissionService', ['$http', function ($http) {
  this.getOne = function (id) {
    return $http.get('/api/submissions/' + id + '/');
  };
  this.getAll = function () {
    return $http.get('/api/submissions/');
  };
}]);
