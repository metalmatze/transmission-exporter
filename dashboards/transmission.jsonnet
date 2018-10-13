local g = import 'grafana-builder/grafana.libsonnet';

{
  grafanaDashboards+:: {
    'transmission.json':
      g.dashboard(
        'Transmission',
      )
      .addMultiTemplate('torrent', 'transmission_torrent_added', 'name')
      .addRow(
        g.row('Downloads')
        .addPanel(
          g.panel('Downloads (bytes/second)') +
          g.queryPanel('max(transmission_torrent_download_bytes{name=~"$torrent"}) by (name) > 0', '{{name}}') +
          g.stack +
          { yaxes: g.yaxes('bytes') },
        )
      )
      .addRow(
        g.row('Upload')
        .addPanel(
          g.panel('Uploads (bytes/second)') +
          g.queryPanel('max(transmission_torrent_upload_bytes{name=~"$torrent"}) by (name) > 0', '{{name}}') +
          g.stack +
          { yaxes: g.yaxes('bytes') },
        )
      )
      .addRow(
        g.row('Torrents')
        .addPanel(
          g.panel('Overview') +
          g.tablePanel([
            'max(transmission_torrent_done{name=~"$torrent"}) by (name)',
            'max(transmission_torrent_ratio{name=~"$torrent"}) by (name)',
            'max(transmission_torrent_leechers{name=~"$torrent"}) by (name)',
            'max(transmission_torrent_seeders{name=~"$torrent"}) by (name)',
          ], {
            'Value #A': { alias: 'Progress', unit: 'percentunit' },
            'Value #B': { alias: 'Ratio' },
            'Value #C': { alias: 'Leechers' },
            'Value #D': { alias: 'Seeders' },
          })
        )
      ),
  },
}
