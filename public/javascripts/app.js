$(document).on('click', '.f-dropdown a', function(e) {
  e.preventDefault();
  $(this).parent().parent().css('left', '-9999999px')
  $.post(this.getAttribute('href'));
})
$(document).foundation();