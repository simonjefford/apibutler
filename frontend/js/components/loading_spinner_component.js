/* global App, Spinner */

App.LoadingSpinnerComponent = Ember.Component.extend({
    spinner: undefined,
    lines: 13,
    length: 20,
    width: 10,
    radius: 30,

    showSpinner: function() {
        var target = this.get('element');
        this.spinner = new Spinner({
            lines: this.get('lines'),
            length: this.get('length'),
            width: this.get('width'),
            radius: this.get('radius')
        });
        this.spinner.spin(target);
    }.on('didInsertElement'),

    teardown: function() {
        if (this.spinner) {
            this.spinner.stop();
        }
    }.on('willDestroyElement')
});
