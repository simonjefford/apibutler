var isEmpty = Ember.isEmpty, keys = Ember.keys;

var ModelMixin = Ember.Mixin.create({
    isClean: function() {
        return isEmpty(keys(this.changedAttributes()));
    }.property().volatile()
});

export default ModelMixin;
