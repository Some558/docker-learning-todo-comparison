using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.EntityFrameworkCore;
using SimpleTodoAPI.Data;
using SimpleTodoAPI.Models;

namespace SimpleTodoAPI.Pages.Todos
{
    public class EditModel : PageModel
    {
        private readonly TodoContext _context;

        public EditModel(TodoContext context)
        {
            _context = context;
        }

        [BindProperty]
        public Todo Todo { get; set; } = default!;

        public async Task<IActionResult> OnGetAsync(int? id)
        {
            if (id == null)
            {
                return NotFound();
            }

            var todo = await _context.Todos.FirstOrDefaultAsync(m => m.Id == id);
            if (todo == null)
            {
                return NotFound();
            }
            Todo = todo;
            return Page();
        }

        public async Task<IActionResult> OnPostAsync()
        {
            if (!ModelState.IsValid)
            {
                return Page();
            }

            _context.Attach(Todo).State = EntityState.Modified;

            try
            {
                await _context.SaveChangesAsync();
            }
            catch (DbUpdateConcurrencyException)
            {
                if (!TodoExists(Todo.Id))
                {
                    return NotFound();
                }
                else
                {
                    throw;
                }
            }

            return RedirectToPage("/Index");
        }

        private bool TodoExists(int id)
        {
            return _context.Todos.Any(e => e.Id == id);
        }
    }
}