using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.EntityFrameworkCore;
using SimpleTodoAPI.Data;
using SimpleTodoAPI.Models;

namespace SimpleTodoAPI.Pages
{
    public class IndexModel : PageModel
    {
        private readonly TodoContext _context;

        public IndexModel(TodoContext context)
        {
            _context = context;
        }

        public IList<Todo> Todos { get; set; } = default!;

        public async Task OnGetAsync()
        {
            Todos = await _context.Todos
                .OrderByDescending(t => t.CreatedAt)
                .ToListAsync();
        }
    }
}