window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">
  
  //Dynamic host update
  spec.host = window.location.host
  //Dynamic externalDocs docs url update
  spec.externalDocs.url = spec.externalDocs.url.replace(/(localhost:3000)/g, window.location.host);
   //Dynamic auth url update
  spec.securityDefinitions.oauth2.tokenUrl = spec.securityDefinitions.oauth2.tokenUrl.replace(/(localhost:3000)/g, window.location.host);
	
  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    spec: spec,
    dom_id: '#swagger-ui',
    deepLinking: true,
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ]
});

  //</editor-fold>
};
