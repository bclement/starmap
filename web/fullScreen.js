function toEquitorial(lon, lat) {
    var ra = (-lon) / 15.0 + 12.0; 
    var decl = lat;
    return new OpenLayers.LonLat(ra, decl);
}
function toStellar(bounds) {
    var rval = new OpenLayers.Bounds();
    rval.extend(toEquitorial(bounds.left, bounds.bottom));
    rval.extend(toEquitorial(bounds.right, bounds.top));
    return rval;
}
var mousePositionCtrl = new OpenLayers.Control.MousePosition({
    prefix: "Equitorial coordinates: ",
    formatOutput: function(lonlat) {
        eq = toEquitorial(lonlat.lon, lonlat.lat);
        return "( " +eq.lon + ", " + eq.lat + " )";
    }
}
);
layer = new OpenLayers.Layer.WMS( "OpenLayers WMS", "/wms",
        {layers: 'stars'} );
layer.getURL = function (bounds) {
    bounds = this.adjustBounds(bounds);
    bounds = toStellar(bounds)

        var imageSize = this.getImageSize();
    var newParams = {};
    // WMS 1.3 introduced axis order
    var reverseAxisOrder = this.reverseAxisOrder();
    newParams.BBOX = this.encodeBBOX ?
        bounds.toBBOX(null, reverseAxisOrder) :
        bounds.toArray(reverseAxisOrder);
    newParams.WIDTH = imageSize.w;
    newParams.HEIGHT = imageSize.h;
    var requestString = this.getFullRequestString(newParams);
    return requestString;
};
var map = new OpenLayers.Map({
    div: "map",
    layers: [layer],
    center: [0, 0],
    zoom: 3
});

info = new OpenLayers.Control.WMSGetFeatureInfo({
    title: 'Identify features by clicking',
     queryVisible: true,
     eventListeners: {
         getfeatureinfo: function(event) {
             var lonlat = map.getLonLatFromPixel(event.xy);
             map.addPopup(new OpenLayers.Popup.FramedCloud(
                     "chicken",lonlat, null, event.text, null,
                     true));
         }
     }
});
info.buildWMSOptions = function(url, layers, clickPosition, format) {
    var layerNames = [], styleNames = [];
    for (var i = 0, len = layers.length; i < len; i++) {
        if (layers[i].params.LAYERS != null) {
            layerNames = layerNames.concat(layers[i].params.LAYERS);
            styleNames = styleNames.concat(info.getStyleNames(layers[i]));
        }
    }
    var firstLayer = layers[0];
    // use the firstLayer's projection if it matches the map projection -
    // this assumes that all layers will be available in this projection
    var projection = info.map.getProjection();
    var layerProj = firstLayer.projection;
    if (layerProj && layerProj.equals(info.map.getProjectionObject())) {
        projection = layerProj.getCode();
    }
    var lonlat = info.map.getExtent();
    var eq = toStellar(lonlat);
    var bbox = eq.toBBOX(null,
            firstLayer.reverseAxisOrder());

    var params = OpenLayers.Util.extend({
        service: "WMS",
        version: firstLayer.params.VERSION,
        request: "GetFeatureInfo",
        exceptions: firstLayer.params.EXCEPTIONS,
        bbox: bbox,
        feature_count: info.maxFeatures,
        height: info.map.getSize().h,
        width: info.map.getSize().w,
        format: format,
        info_format: firstLayer.params.INFO_FORMAT || info.infoFormat
    }, (parseFloat(firstLayer.params.VERSION) >= 1.3) ?
    {
        crs: projection,
        i: parseInt(clickPosition.x),
        j: parseInt(clickPosition.y)
    } :
    {
        srs: projection,
        x: parseInt(clickPosition.x),
        y: parseInt(clickPosition.y)
    }
    );
    if (layerNames.length != 0) {
        params = OpenLayers.Util.extend({
            layers: layerNames,
            query_layers: layerNames,
            styles: styleNames
        }, params);
    }
    OpenLayers.Util.applyDefaults(params, info.vendorParams);
    return {
        url: url,
            params: OpenLayers.Util.upperCaseObject(params),
            callback: function(request) {
                info.handleResponse(clickPosition, request, url);
            },
            scope: info
    };
};
map.addControl(info);
info.activate();
map.addControl(new OpenLayers.Control.LayerSwitcher());
map.addControl(mousePositionCtrl);
