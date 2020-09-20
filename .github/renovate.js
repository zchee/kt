process.env.GO111MODULE = 'on';
process.env.GOPROXY = 'https://proxy.golang.org,direct';
process.env.GOSUMDB = 'sum.golang.org';

module.exports = {
  extends: [
    'config:base'
  ],
  golang: {
    ignoreDeps: [
      'k8s.io/api',
      'k8s.io/apimachinery',
      'k8s.io/client-go',
      'sigs.k8s.io/controller-runtime',
    ],
    postUpdateOptions: ['gomodTidy'],
  },
  reviewers: [
    'zchee'
  ],
  automerge: true,
  major: {
    automerge: false
  },
  rebaseWhen: 'behind-base-branch',
  labels: ['automerge'],
  packageRules: [
    {
      datasources: 'go',
      managers: ['gomod'],
      updateTypes: ['pin', 'digest'],
      versioning: 'semver'
    },
  ],
  timezone: 'Asia/Tokyo',
};
