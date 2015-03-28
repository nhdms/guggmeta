'use strict';

var app = angular.module('guggmeta', ['ngRoute', 'ngAnimate', 'ui.bootstrap.pagination', 'angular-loading-bar']);

app.run(['$window', '$rootScope', function ($window, $rootScope) {
  $window.onload = function () {
    try {
      console.log("Welcome!\nYou can find the source code of guggmeta in https://github.com/sevein/guggmeta.\nPlease, send me your feedback!");
    } catch (e) {}
  };
  $rootScope.home = false;
  $rootScope.$on('$routeChangeSuccess', function (event, current, previous) {
    $rootScope.home = current.$$route.originalPath === '/';
  });
  $rootScope.top = function () {
    $window.scrollTo(0, 0);
  };
}]);

app.config(['$routeProvider', '$locationProvider', function ($routeProvider, $locationProvider) {
  $routeProvider
    .when('/', {
      templateUrl: '/assets/partials/home.tmpl.html',
      controller: 'HomeCtrl'
    })
    .when('/analytics', {
      templateUrl: '/assets/partials/analytics-list.tmpl.html',
      controller: 'AnalyticsListCtrl',
      resolve: {
        response: ['SubmissionService', function (SubmissionService) {
          return SubmissionService.getAnalytics();
        }]
      }
    })
    .when('/submissions', {
      templateUrl: '/assets/partials/submission-list.tmpl.html',
      controller: 'SubmissionListCtrl',
      resolve: {
        response: ['$location', 'SubmissionService', function ($location, SubmissionService) {
          var query = '';
          var params = $location.search();
          if (params.hasOwnProperty('q')) {
            query = params.q;
          }
          return SubmissionService.search(query);
        }]
      }
    })
    .when('/submissions/:id', {
      templateUrl: '/assets/partials/submission-detail.tmpl.html',
      controller: 'SubmissionDetailCtrl',
      resolve: {
        response: ['$route', 'SubmissionService', function ($route, SubmissionService) {
          return SubmissionService.getOne($route.current.params.id);
        }]
      }
    })
    .otherwise('/');
  $locationProvider.html5Mode({
    enabled: true,
    requireBase: false
  });
}]);

app.controller('HeaderCtrl', ['$scope', '$location', '$window', function ($scope, $location, $window) {
  $scope.form = {};
  $scope.search = function () {
    if (angular.isUndefined($scope.query)) {
      return;
    }
    if (!$scope.query.length) {
      $location.path('/submissions');
      return;
    }
    if (/GH-\d/.test($scope.query)) {
      $location.path('/submissions/' + $scope.query);
      return;
    }
    $location.path('/submissions').search({ 'q': $scope.query });
  };
  $scope.$on('$routeChangeError', function () {
    $window.alert('Not found!');
    delete $scope.query;
  });
  $scope.$on('$routeChangeSuccess', function (event, current, previous) {
    if ($location.path() !== '/submissions') {
      delete $scope.query;
    } else {
      $scope.query = $location.search()['q'];
    }
  });
}]);

app.controller('HomepageTopCtrl', ['$scope', '$location', '$window', function ($scope, $location, $window) {
  $scope.form = {};
  $scope.search = function () {
    if (angular.isUndefined($scope.query)) {
      return;
    }
    if (!$scope.query.length) {
      $location.path('/submissions');
      return;
    }
    if (/GH-\d/.test($scope.query)) {
      $location.path('/submissions/' + $scope.query);
      return;
    }
    $location.path('/submissions').search({ 'q': $scope.query });
  };
}]);

app.controller('HomeCtrl', [function () {

}]);

app.controller('AnalyticsListCtrl', ['$scope', 'response', function ($scope, response) {
  $scope.analytics = response.data;
}]);

app.controller('SubmissionListCtrl', ['$scope', '$location', '$sce', 'response', 'SubmissionService', function ($scope, $location, $sce, response, SubmissionService) {
  var params = $location.search();
  $scope.query = params.hasOwnProperty('q') ? params.q : undefined;
  $scope.pager = { currentPage: 1, itemsPerPage: 10, maxSize: 5 }
  var populate = function (resp) {
    $scope.totalItems = resp.data.total;
    $scope.submissions = [];
    var length = resp.data.results.length;
    for (var i = 0; i < length; i++) {
      var tuple = resp.data.results[i];
      if (angular.isDefined(tuple.highlight) && tuple.highlight !== null) {
        tuple.highlight = $sce.trustAsHtml('[...] ' + tuple.highlight['pdfs.content'][0])
      }
      $scope.submissions.push(tuple);
    }
  };
  populate(response);
  $scope.pageChanged = function () {
    SubmissionService.search($scope.query, $scope.pager.currentPage).then(function (response) {
      populate(response);
      $scope.top();
    });
  };
}]);

app.controller('SubmissionDetailCtrl', ['$rootScope', '$scope', '$routeParams', 'response', function ($rootScope, $scope, $routeParams, response) {
  $scope.id = $routeParams.id;
  $scope.submission = response.data;
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
  $scope.open = function (item) {
    item.preview = true;
    $rootScope.popup = true;
  }
  $scope.close = function (item) {
    item.preview = false;
    $rootScope.popup = false;
  }
}]);

app.service('SubmissionService', ['$http', function ($http) {
  this.getOne = function (id) {
    return $http.get('/api/submissions/' + id + '/');
  };
  this.getAll = function () {
    return this.search();
  };
  this.getAnalytics = function () {
    return $http.get('/api/submissions/analytics/');
  };
  this.search = function (query, page) {
    var params = {};
    if (angular.isDefined(query)) {
      params.q = query;
    }
    if (angular.isDefined(page)) {
      params.p = page;
    }
    return $http.get('/api/submissions/', {
      params: params
    });
  };
}]);
