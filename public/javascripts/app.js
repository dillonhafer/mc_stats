(function() {
  function setNames() {
    $.get("/players", function(players) {
      for(var i in players) {
        player = players[i]
        var dom_id = player.uuid
        $('#uuid-'+dom_id).html(player.name)
      }
    })
  }

  function createTitle(uuid) {
    title = document.createElement('li')
    title.setAttribute('class', 'title')
    title.setAttribute('id', 'uuid-'+ uuid)
    return title
  }

  function createItem(label, value) {
    var name = label.split('.').slice(-1)[0].replace(/-/g, ' ')
    li = document.createElement('li')
    li.setAttribute('data-'+label, value)
    li.setAttribute('class', 'bullet-item')
    li.innerHTML = name+": "+value
    return li
  }

  function createPlayer(player) {
    var li = document.createElement('li')
    var ul = document.createElement('ul')
    ul.setAttribute('class', 'pricing-table')
    playerTitle = createTitle(player.UUID)
    ul.appendChild(playerTitle)

    for (var key in player.data) {
      if (player.data.hasOwnProperty(key)) {
        item = createItem(key, player.data[key])
        ul.appendChild(item)
      }
    }

    li.appendChild(ul)
    return li
  }

  $(document).foundation();

  $.get("/stats", function(players) {
    for(var i in players) {
      player = createPlayer(players[i]);
      $('.main').append(player)
    }
    setNames()
  })
}).call(this)
