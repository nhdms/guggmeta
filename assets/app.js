'use strict';

var app = angular.module('guggmeta', ['ngRoute', 'ngAria', 'ngAnimate']);

app.run(['$window', '$rootScope', function ($window, $rootScope) {
  $window.onload = function () {
    try {
      console.log("Welcome!\nYou can find the source code of guggmeta in https://github.com/sevein/guggmeta.\nPlease, send me your feedback!");
    } catch (e) {}
  };
  $rootScope.$on('$routeChangeStart', function (event, next, current) {
    $rootScope.home = next.$$route.originalPath === '/';
  });
}]);

app.config(['$routeProvider', '$locationProvider', function ($routeProvider, $locationProvider) {
  $routeProvider
    .when('/', {
      templateUrl: '/assets/partials/home.tmpl.html',
      controller: 'HomeCtrl'
    })
    .when('/analytics', {
      templateUrl: '/assets/partials/analytics-list.tmpl.html',
      controller: 'AnalyticsListCtrl'
    })
    .when('/submissions', {
      templateUrl: '/assets/partials/submission-list.tmpl.html',
      controller: 'SubmissionListCtrl'
    })
    .when('/submissions/:id', {
      templateUrl: '/assets/partials/submission-detail.tmpl.html',
      controller: 'SubmissionDetailCtrl'
    })
    .otherwise('/');
  $locationProvider.html5Mode({
    enabled: true,
    requireBase: false
  });
}]);

app.controller('HomeCtrl', ['$scope', function ($scope) {

}]);

app.controller('AnalyticsListCtrl', ['$scope', 'SubmissionService', function ($scope, SubmissionService) {
  SubmissionService.getAnalytics().then(function (response) {
    $scope.analytics = response.data;
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
  // Properties "file_name" and "author" are ignored or handled manually
  $scope.submissionFields = [
    { 'label': 'Creation date', 'property': 'creation_date' },
    { 'label': 'Creator', 'property': 'creator' },
    { 'label': 'Encrypted', 'property': 'encrypted' },
    { 'label': 'Size', 'property': 'file_size', 'suffix': ' bytes' },
    { 'label': 'Form', 'property': 'form' },
    { 'label': 'JavaScript', 'property': 'javascript' },
    { 'label': 'Keywords', 'property': 'keywords' },
    { 'label': 'Modification date', 'property': 'mod_date' },
    { 'label': 'Optimized', 'property': 'optimized' },
    { 'label': 'Page rot', 'property': 'page_rot' },
    { 'label': 'Page size', 'property': 'page_size' },
    { 'label': 'Pages', 'property': 'pages' },
    { 'label': 'PDF version', 'property': 'pdf_version' },
    { 'label': 'Producer', 'property': 'producer' },
    { 'label': 'Subject', 'property': 'subject' },
    { 'label': 'Suspects', 'property': 'suspects' },
    { 'label': 'Tagged', 'property': 'tagged' },
    { 'label': 'Title', 'property': 'title' },
    { 'label': 'UserProperties', 'property': 'user_properties' },
  ];
}]);

app.service('SubmissionService', ['$http', function ($http) {
  this.getOne = function (id) {
    return $http.get('/api/submissions/' + id + '/');
  };
  this.getAll = function () {
    return $http.get('/api/submissions/');
  };
  this.getAnalytics = function () {
    return $http.get('/api/submissions/analytics/');
  };
}]);
