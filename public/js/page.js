$(function(){

  function init() {
    initUI()
    initMarked()
    convert()
  }

  function initUI() {
    $('.alert').hide()
    $('.modal').hide()

    $('#createModalInput').val(syaro.wikiPath)
    $('#renameModalInput').val(syaro.wikiPath)

    $('#createModalButton').on('click', function() {
      var name = $('#createModalInput').val()
      if (name === "") {
        $('#createErrorAlert').html('<strong>Error</strong> Please fill brank form.')
        $('#createErrorAlert').show()
        return
      }
      if (name[0] !== '/') { name = '/' + name }

      var reqUrl = location.href.replace(syaro.wikiPath,
          encodeURIComponent(name).replace(/%2F/g, '/'))

      var req = new XMLHttpRequest()
      req.open('GET', reqUrl + '?action=create')

      req.onreadystatechange = function() {
        if (req.readyState === 4) {
          $('#createModalButton').button('reset')

          switch (req.status) {
          case 200:
            // redirect to editor
            location.href = reqUrl + '?view=editor'
            break

          default:
            // show error alert
            $('#createErrorAlert').html('<strong>Error</strong> ' + req.responseText)
            $('#createErrorAlert').show()
            break
          }
        }
      }

      req.send()
      $('#createModalButton').button('loading')
    })

    $('#renameModalButton').on('click', function() {
      var name = $('#renameModalInput').val()
      if (name === "") {
        $('#renameErrorAlert').html('<strong>Error</strong> Please fill brank form.')
        $('#renameErrorAlert').show()
        return
      }
      if (name[0] !== '/') { name = '/' + name }

      var reqUrl = location.href.replace(syaro.wikiPath,
          encodeURIComponent(name).replace(/%2F/g, '/'))

      var req = new XMLHttpRequest()
      req.open('GET', reqUrl + '?action=rename&oldpath=' + encodeURIComponent(wikiName))

      req.onreadystatechange = function() {
        if (req.readyState === 4) {
          $('#renameModalButton').button('reset')

          switch (req.status) {
          case 200:
            // redirect to page
            location.href = reqUrl
            break

          default:
            // show error alert
            $('#renameErrorAlert').html('<strong>Error</strong> ' + req.responseText)
            $('#renameErrorAlert').show()
            break
          }
        }
      }

      req.send()
      $('#renameModalButton').button('loading')
    })

    $('#deleteModalButton').on('click', function() {
      var reqUrl = location.href

      var req = new XMLHttpRequest()
      req.open('GET', reqUrl + '?action=delete')

      req.onreadystatechange = function() {
        if (req.readyState === 4) {
          $('#deleteModalButton').button('reset')

          switch (req.status) {
          case 200:
            $('#deleteErrorAlert').hide()
            // show success alert
            $('#deleteSuccessAlert').show()
            break

          default:
            // show error alert
            $('#deleteErrorAlert').html('<strong>Error</strong> ' + req.responseText)
            $('#deleteErrorAlert').show()
            break
          }
        }
      }

      req.send()
      $('#deleteModalButton').button('loading')
    })
  }

  function initMarked() {
    marked.setOptions({
      gfm: true,
      tables: true,
      pedantic: false,
      sanitize: false,
      smartLists: true,
      smartypants: false,
      langPrefix: 'lang-'
      // highlight: function (code) {
      //   return hljs.highlightAuto(code).value
      // }
    })
  }

  function convert() {
    var $mainMd = $('.syaro-main > .markdown')
    var $sidebarMd = $('.syaro-sidebar > .markdown')

    $mainMd.html( marked($mainMd.html()) )
    $sidebarMd.html( marked($sidebarMd.html()) )
  }

  init()

})
