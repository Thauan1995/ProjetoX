$('#nova-publicacao').on('submit', criarPublicacao);

function criarPublicacao(evento) {
    evento.preventDefault();

    $.ajax({
        url: "/web/publicacoes",
        method: "POST",
        data: {
            titulo: $('#titulo').val(),
            conteudo: $('#conteudo').val(),
        }
    }).done(function() {
        window.location = "/web/home";
    }).fail(function() {
        alert("Erro ao criar a publicação!");
    });
}